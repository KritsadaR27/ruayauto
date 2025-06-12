#!/bin/bash

# Microservices Management Script
# FacebookConnect + ChatBotCore Architecture

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if Docker is running
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker and try again."
        exit 1
    fi
}

# Function to build all services
build() {
    print_status "Building microservices..."
    docker-compose -f docker-compose.microservices.yml build --no-cache
    print_success "All services built successfully!"
}

# Function to start all services
start() {
    print_status "Starting microservices..."
    docker-compose -f docker-compose.microservices.yml up -d
    
    print_status "Waiting for services to be healthy..."
    sleep 15
    
    # Check service health
    check_health
}

# Function to stop all services
stop() {
    print_status "Stopping microservices..."
    docker-compose -f docker-compose.microservices.yml down
    print_success "All services stopped!"
}

# Function to restart services
restart() {
    print_status "Restarting microservices..."
    docker-compose -f docker-compose.microservices.yml down
    docker-compose -f docker-compose.microservices.yml up -d
    print_success "All services restarted!"
}

# Function to show logs
logs() {
    if [ -n "$2" ]; then
        docker-compose -f docker-compose.microservices.yml logs -f "$2"
    else
        docker-compose -f docker-compose.microservices.yml logs -f
    fi
}

# Function to check service health
check_health() {
    print_status "Checking service health..."
    
    # Check ChatBotCore
    if curl -s http://localhost:8086/health > /dev/null 2>&1; then
        print_success "✅ ChatBotCore is healthy (port 8086)"
    else
        print_error "❌ ChatBotCore is unhealthy"
    fi
    
    # Check FacebookConnect
    if curl -s http://localhost:8085/health > /dev/null 2>&1; then
        print_success "✅ FacebookConnect is healthy (port 8085)"
    else
        print_error "❌ FacebookConnect is unhealthy"
    fi
    
    # Check Database
    if docker-compose -f docker-compose.microservices.yml exec -T postgres pg_isready -U chatbot_user -d chatbot_mvp > /dev/null 2>&1; then
        print_success "✅ PostgreSQL is healthy"
    else
        print_error "❌ PostgreSQL is unhealthy"
    fi
}

# Function to show status
status() {
    print_status "Service Status:"
    docker-compose -f docker-compose.microservices.yml ps
    
    echo ""
    check_health
    
    echo ""
    print_status "Service URLs:"
    echo "  🤖 ChatBotCore (Internal):   http://localhost:8086"
    echo "  📘 FacebookConnect (External): http://localhost:8085"
    echo "  🗄️  PostgreSQL:              localhost:55433"
    echo "  🌐 Frontend (if enabled):    http://localhost:3000"
}

# Function to run tests
test() {
    print_status "Running service integration tests..."
    
    # Test ChatBotCore health
    print_status "Testing ChatBotCore..."
    curl -f http://localhost:8086/health || {
        print_error "ChatBotCore health check failed"
        return 1
    }
    
    # Test FacebookConnect health
    print_status "Testing FacebookConnect..."
    curl -f http://localhost:8085/health || {
        print_error "FacebookConnect health check failed"
        return 1
    }
    
    # Test service communication
    print_status "Testing service communication..."
    curl -f -X POST http://localhost:8086/api/v1/messages/process \
        -H "Content-Type: application/json" \
        -d '{"content":"test","sender_id":"test","conversation_id":"test","timestamp":"2024-01-01T00:00:00Z","message_type":"text"}' || {
        print_error "Service communication test failed"
        return 1
    }
    
    print_success "All tests passed!"
}

# Function to clean up everything
clean() {
    print_warning "This will remove all containers, images, and volumes!"
    read -p "Are you sure? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        print_status "Cleaning up..."
        docker-compose -f docker-compose.microservices.yml down -v --rmi all
        docker system prune -f
        print_success "Cleanup completed!"
    else
        print_status "Cleanup cancelled."
    fi
}

# Function to show help
help() {
    echo "Microservices Management Script"
    echo "FacebookConnect (External) + ChatBotCore (Internal) Architecture"
    echo ""
    echo "Usage: $0 <command> [options]"
    echo ""
    echo "Commands:"
    echo "  build       Build all Docker images"
    echo "  start       Start all services"
    echo "  stop        Stop all services"
    echo "  restart     Restart all services"
    echo "  logs        Show logs (use 'logs <service>' for specific service)"
    echo "  status      Show service status and health"
    echo "  test        Run integration tests"
    echo "  clean       Remove all containers, images, and volumes"
    echo "  help        Show this help message"
    echo ""
    echo "Service Names:"
    echo "  - chatbot-core       Internal business logic service"
    echo "  - facebook-connect   External Facebook integration service"
    echo "  - postgres           PostgreSQL database"
    echo "  - frontend           Next.js frontend (optional)"
    echo ""
    echo "Examples:"
    echo "  $0 start                    # Start all services"
    echo "  $0 logs chatbot-core        # Show ChatBotCore logs"
    echo "  $0 logs facebook-connect    # Show FacebookConnect logs"
    echo "  $0 test                     # Run integration tests"
}

# Main script logic
main() {
    check_docker
    
    case "${1:-help}" in
        build)
            build
            ;;
        start)
            start
            ;;
        stop)
            stop
            ;;
        restart)
            restart
            ;;
        logs)
            logs "$@"
            ;;
        status)
            status
            ;;
        test)
            test
            ;;
        clean)
            clean
            ;;
        help|--help|-h)
            help
            ;;
        *)
            print_error "Unknown command: $1"
            help
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"

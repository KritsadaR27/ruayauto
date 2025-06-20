export type Response = {
  text: string
  image?: string // URL or base64 image for this specific response
}

export type ConnectedPage = {
  id: string
  name: string
  pageId: string
  connected: boolean
  enabled: boolean
}

export type Pair = {
  id?: string // Add unique ID for each pair
  title?: string // Custom title for each pair 
  hasManuallyEditedTitle?: boolean // Track if user manually edited title
  keywords: string[]
  responses: Response[] // Change from string[] to Response[]
  // Individual settings for each pair
  hideCommentsAfterReply?: boolean
  enableInboxIntegration?: boolean
  inboxResponse?: string
  inboxImage?: string
  // New design properties
  enabled?: boolean // Toggle for enabling/disabling the rule
  expanded?: boolean // Toggle for expanding/collapsing the card
  // Multi-page support
  selectedPages?: string[] // Array of page IDs this rule applies to
}

export type FilterSettings = {
  skipMentions: boolean  // Skip comments with @ mentions
  skipStickers: boolean  // Skip comments with emoji/stickers
}

export type FallbackSettings = {
  enabled: boolean
  responses: Response[]
  // Multi-page support for fallback
  selectedPages?: string[]
  hideAfterReply?: boolean
  sendToInbox?: boolean
  inboxMessage?: string
  inboxImage?: string
}

export type Settings = {
  defaultResponses: string[]
  enableDefault: boolean
  noTag: boolean
  noSticker: boolean
  hideCommentsAfterReply: boolean
  enableInboxIntegration: boolean
  inboxResponse: string
  inboxImage?: string
  // New settings
  filterSettings: FilterSettings
  fallbackSettings: FallbackSettings
}

export type KeywordData = {
  pairs: Pair[]
} & Settings

// Additional types for the ruayAutoMsg system
export type Keyword = {
  id: number
  keyword: string
  response: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export type KeywordResponse = {
  pairs: Keyword[]
  total?: number
  success?: boolean
}

export type ApiResponse<T = any> = {
  success: boolean
  data?: T
  error?: string
  message?: string
  timestamp?: string
}

// Refactor: rename to rule.ts, change all keyword(s) to rule(s) in types and comments

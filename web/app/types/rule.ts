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

export type RulePair = {
  id?: string // Add unique ID for each pair
  title?: string // Custom title for each pair 
  hasManuallyEditedTitle?: boolean // Track if user manually edited title
  rules: string[]
  responses: Response[]
  hideCommentsAfterReply?: boolean
  enableInboxIntegration?: boolean
  inboxResponse?: string
  inboxImage?: string
  enabled?: boolean
  expanded?: boolean
  selectedPages?: string[]
}

export type FilterSettings = {
  skipMentions: boolean
  skipStickers: boolean
}

export type FallbackSettings = {
  enabled: boolean
  responses: Response[]
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
  filterSettings: FilterSettings
  fallbackSettings: FallbackSettings
}

export type RuleData = {
  rules: RulePair[]
} & Settings

export type Rule = {
  id: number
  keywords: string[]
  response: string
  is_active: boolean
  priority: number
  match_type: string
  created_at: string
  updated_at: string
  rule_name: string
  hide_after_reply: boolean
  send_to_inbox: boolean
  inbox_message: string
  inbox_image: string
  created_by?: number | null
}

export type RuleResponse = {
  pairs: Rule[]
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

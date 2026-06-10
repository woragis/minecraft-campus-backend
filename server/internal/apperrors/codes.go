package apperrors

const (
	CodeInternal              = "INTERNAL_ERROR"
	CodeDependencyUnavailable = "DEPENDENCY_UNAVAILABLE"
)

// Plugin auth
const (
	CodePluginAuthV1HandlerKeyMissing = "PLUGIN_AUTH_V1_HANDLER_KEY_MISSING"
	MsgPluginAuthV1HandlerKeyMissing  = "Plugin API key is required."

	CodePluginAuthV1HandlerKeyInvalid = "PLUGIN_AUTH_V1_HANDLER_KEY_INVALID"
	MsgPluginAuthV1HandlerKeyInvalid  = "Plugin API key is invalid."
)

// GET /ready
const (
	CodeReadyGetHandlerSqlGetterFailed = "READY_GET_HANDLER_SQL_GETTER_FAILED"
	MsgReadyGetHandlerSqlGetterFailed  = "Database connection handle is not available."

	CodeReadyGetHandlerDatabasePingFailed = "READY_GET_HANDLER_DATABASE_PING_FAILED"
	MsgReadyGetHandlerDatabasePingFailed  = "Database is not reachable."
)

// GET /v1/internal/whitelist/{minecraftUuid}
const (
	CodeWhitelistGetV1HandlerUUIDInvalid = "WHITELIST_GET_V1_HANDLER_UUID_INVALID"
	MsgWhitelistGetV1HandlerUUIDInvalid  = "minecraftUuid must be a valid UUID."

	CodeWhitelistGetV1ServiceCheckFailed = "WHITELIST_GET_V1_SERVICE_CHECK_FAILED"
	MsgWhitelistGetV1ServiceCheckFailed  = "Could not check whitelist."
)

// POST /v1/internal/players/upsert
const (
	CodePlayerUpsertV1HandlerBodyInvalid = "PLAYER_UPSERT_V1_HANDLER_BODY_INVALID"
	MsgPlayerUpsertV1HandlerBodyInvalid  = "Request body must include minecraftUuid, username, and serverSlug."

	CodePlayerUpsertV1ServiceFailed = "PLAYER_UPSERT_V1_SERVICE_FAILED"
	MsgPlayerUpsertV1ServiceFailed  = "Could not upsert player."
)

// POST /v1/internal/invites
const (
	CodeInvitePostInternalV1HandlerBodyInvalid = "INVITE_POST_INTERNAL_V1_HANDLER_BODY_INVALID"
	MsgInvitePostInternalV1HandlerBodyInvalid  = "Request body must include sponsorUuid and targetUsername."

	CodeInvitePostInternalV1ServiceSponsorNotFound = "INVITE_POST_INTERNAL_V1_SERVICE_SPONSOR_NOT_FOUND"
	MsgInvitePostInternalV1ServiceSponsorNotFound  = "Sponsor player not found."

	CodeInvitePostInternalV1ServiceSponsorBanned = "INVITE_POST_INTERNAL_V1_SERVICE_SPONSOR_BANNED"
	MsgInvitePostInternalV1ServiceSponsorBanned  = "Banned players cannot send invites."

	CodeInvitePostInternalV1ServicePendingExists = "INVITE_POST_INTERNAL_V1_SERVICE_PENDING_EXISTS"
	MsgInvitePostInternalV1ServicePendingExists  = "A pending invite already exists for this username."

	CodeInvitePostInternalV1ServiceCreateFailed = "INVITE_POST_INTERNAL_V1_SERVICE_CREATE_FAILED"
	MsgInvitePostInternalV1ServiceCreateFailed  = "Could not create invite."
)

// GET /v1/players/{id}
const (
	CodePlayerGetV1HandlerIDInvalid = "PLAYER_GET_V1_HANDLER_ID_INVALID"
	MsgPlayerGetV1HandlerIDInvalid  = "Player id must be a valid UUID."

	CodePlayerGetV1ServiceNotFound = "PLAYER_GET_V1_SERVICE_NOT_FOUND"
	MsgPlayerGetV1ServiceNotFound  = "Player not found."

	CodePlayerGetV1ServiceLoadFailed = "PLAYER_GET_V1_SERVICE_LOAD_FAILED"
	MsgPlayerGetV1ServiceLoadFailed  = "Could not load player."
)

// GET /v1/players/minecraft/{minecraftUuid}
const (
	CodePlayerGetMinecraftV1HandlerUUIDInvalid = "PLAYER_GET_MINECRAFT_V1_HANDLER_UUID_INVALID"
	MsgPlayerGetMinecraftV1HandlerUUIDInvalid  = "minecraftUuid must be a valid UUID."

	CodePlayerGetMinecraftV1ServiceNotFound = "PLAYER_GET_MINECRAFT_V1_SERVICE_NOT_FOUND"
	MsgPlayerGetMinecraftV1ServiceNotFound  = "Player not found."

	CodePlayerGetMinecraftV1ServiceLoadFailed = "PLAYER_GET_MINECRAFT_V1_SERVICE_LOAD_FAILED"
	MsgPlayerGetMinecraftV1ServiceLoadFailed  = "Could not load player."
)

// GET /v1/players/{id}/invites
const (
	CodePlayerInvitesListV1HandlerIDInvalid = "PLAYER_INVITES_LIST_V1_HANDLER_ID_INVALID"
	MsgPlayerInvitesListV1HandlerIDInvalid  = "Player id must be a valid UUID."

	CodePlayerInvitesListV1ServiceNotFound = "PLAYER_INVITES_LIST_V1_SERVICE_NOT_FOUND"
	MsgPlayerInvitesListV1ServiceNotFound  = "Player not found."

	CodePlayerInvitesListV1ServiceListFailed = "PLAYER_INVITES_LIST_V1_SERVICE_LIST_FAILED"
	MsgPlayerInvitesListV1ServiceListFailed  = "Could not list invites."
)

// GET /v1/invites/{code}
const (
	CodeInviteGetV1HandlerCodeEmpty = "INVITE_GET_V1_HANDLER_CODE_EMPTY"
	MsgInviteGetV1HandlerCodeEmpty  = "Invite code is required."

	CodeInviteGetV1ServiceNotFound = "INVITE_GET_V1_SERVICE_NOT_FOUND"
	MsgInviteGetV1ServiceNotFound  = "Invite not found."

	CodeInviteGetV1ServiceLoadFailed = "INVITE_GET_V1_SERVICE_LOAD_FAILED"
	MsgInviteGetV1ServiceLoadFailed  = "Could not load invite."
)

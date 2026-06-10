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

	CodeInvitePostInternalV1ServiceSponsorProbation = "INVITE_POST_INTERNAL_V1_SERVICE_SPONSOR_PROBATION"
	MsgInvitePostInternalV1ServiceSponsorProbation  = "Players on probation cannot send invites."

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

// Guilds
const (
	CodeGuildPostV1HandlerBodyInvalid = "GUILD_POST_V1_HANDLER_BODY_INVALID"
	MsgGuildPostV1HandlerBodyInvalid  = "Request body must include leaderUuid and name."

	CodeGuildPostV1ServiceNameInvalid = "GUILD_POST_V1_SERVICE_NAME_INVALID"
	MsgGuildPostV1ServiceNameInvalid  = "Guild name must be between 1 and 64 characters."

	CodeGuildPostV1ServiceLeaderNotFound = "GUILD_POST_V1_SERVICE_LEADER_NOT_FOUND"
	MsgGuildPostV1ServiceLeaderNotFound  = "Leader player not found."

	CodeGuildPostV1ServiceLeaderBanned = "GUILD_POST_V1_SERVICE_LEADER_BANNED"
	MsgGuildPostV1ServiceLeaderBanned  = "Banned players cannot create guilds."

	CodeGuildPostV1ServiceLeaderProbation = "GUILD_POST_V1_SERVICE_LEADER_PROBATION"
	MsgGuildPostV1ServiceLeaderProbation  = "Players on probation cannot create guilds."

	CodeGuildPostV1ServiceAlreadyMember = "GUILD_POST_V1_SERVICE_ALREADY_MEMBER"
	MsgGuildPostV1ServiceAlreadyMember  = "Player already belongs to a guild."

	CodeGuildPostV1ServiceCreateFailed = "GUILD_POST_V1_SERVICE_CREATE_FAILED"
	MsgGuildPostV1ServiceCreateFailed  = "Could not create guild."

	CodeGuildGetV1HandlerIDInvalid = "GUILD_GET_V1_HANDLER_ID_INVALID"
	MsgGuildGetV1HandlerIDInvalid  = "Guild id must be a valid UUID."

	CodeGuildGetV1ServiceNotFound = "GUILD_GET_V1_SERVICE_NOT_FOUND"
	MsgGuildGetV1ServiceNotFound  = "Guild not found."

	CodeGuildGetV1ServiceLoadFailed = "GUILD_GET_V1_SERVICE_LOAD_FAILED"
	MsgGuildGetV1ServiceLoadFailed  = "Could not load guild."

	CodeGuildListV1ServiceListFailed = "GUILD_LIST_V1_SERVICE_LIST_FAILED"
	MsgGuildListV1ServiceListFailed  = "Could not list guilds."

	CodeGuildMembersV1ServiceGuildNotFound = "GUILD_MEMBERS_V1_SERVICE_GUILD_NOT_FOUND"
	MsgGuildMembersV1ServiceGuildNotFound  = "Guild not found."

	CodeGuildMembersV1ServiceListFailed = "GUILD_MEMBERS_V1_SERVICE_LIST_FAILED"
	MsgGuildMembersV1ServiceListFailed  = "Could not list guild members."

	CodeGuildJoinV1ServicePlayerNotFound = "GUILD_JOIN_V1_SERVICE_PLAYER_NOT_FOUND"
	MsgGuildJoinV1ServicePlayerNotFound  = "Player not found."

	CodeGuildJoinV1ServiceGuildNotFound = "GUILD_JOIN_V1_SERVICE_GUILD_NOT_FOUND"
	MsgGuildJoinV1ServiceGuildNotFound  = "Guild not found."

	CodeGuildJoinV1ServiceProbation = "GUILD_JOIN_V1_SERVICE_PROBATION"
	MsgGuildJoinV1ServiceProbation  = "Players on probation cannot join guilds."

	CodeGuildJoinV1ServiceAlreadyMember = "GUILD_JOIN_V1_SERVICE_ALREADY_MEMBER"
	MsgGuildJoinV1ServiceAlreadyMember  = "Player already belongs to a guild."

	CodeGuildJoinV1ServiceFailed = "GUILD_JOIN_V1_SERVICE_FAILED"
	MsgGuildJoinV1ServiceFailed  = "Could not join guild."

	CodeGuildLeaveV1ServicePlayerNotFound = "GUILD_LEAVE_V1_SERVICE_PLAYER_NOT_FOUND"
	MsgGuildLeaveV1ServicePlayerNotFound  = "Player not found."

	CodeGuildLeaveV1ServiceGuildNotFound = "GUILD_LEAVE_V1_SERVICE_GUILD_NOT_FOUND"
	MsgGuildLeaveV1ServiceGuildNotFound  = "Guild not found."

	CodeGuildLeaveV1ServiceLeaderCannotLeave = "GUILD_LEAVE_V1_SERVICE_LEADER_CANNOT_LEAVE"
	MsgGuildLeaveV1ServiceLeaderCannotLeave  = "Guild leader cannot leave; transfer leadership first."

	CodeGuildLeaveV1ServiceNotMember = "GUILD_LEAVE_V1_SERVICE_NOT_MEMBER"
	MsgGuildLeaveV1ServiceNotMember  = "Player is not a member of this guild."

	CodeGuildLeaveV1ServiceFailed = "GUILD_LEAVE_V1_SERVICE_FAILED"
	MsgGuildLeaveV1ServiceFailed  = "Could not leave guild."
)

// Trust
const (
	CodeTrustEventV1HandlerBodyInvalid = "TRUST_EVENT_V1_HANDLER_BODY_INVALID"
	MsgTrustEventV1HandlerBodyInvalid  = "Request body must include playerId and eventType."

	CodeTrustEventV1ServicePlayerNotFound = "TRUST_EVENT_V1_SERVICE_PLAYER_NOT_FOUND"
	MsgTrustEventV1ServicePlayerNotFound  = "Player not found."

	CodeTrustEventV1ServiceRecordFailed = "TRUST_EVENT_V1_SERVICE_RECORD_FAILED"
	MsgTrustEventV1ServiceRecordFailed  = "Could not record trust event."

	CodeTrustListV1ServicePlayerNotFound = "TRUST_LIST_V1_SERVICE_PLAYER_NOT_FOUND"
	MsgTrustListV1ServicePlayerNotFound  = "Player not found."

	CodeTrustListV1ServiceListFailed = "TRUST_LIST_V1_SERVICE_LIST_FAILED"
	MsgTrustListV1ServiceListFailed  = "Could not list trust events."

	CodeTrustTreeV1ServicePlayerNotFound = "TRUST_TREE_V1_SERVICE_PLAYER_NOT_FOUND"
	MsgTrustTreeV1ServicePlayerNotFound  = "Player not found."

	CodeTrustTreeV1ServiceLoadFailed = "TRUST_TREE_V1_SERVICE_LOAD_FAILED"
	MsgTrustTreeV1ServiceLoadFailed  = "Could not load sponsor tree."
)

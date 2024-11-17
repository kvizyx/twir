package commands

import (
	"github.com/samber/lo"
	seventv "github.com/satont/twir/apps/parser/internal/commands/7tv"
	channel_game "github.com/satont/twir/apps/parser/internal/commands/channel/game"
	channel_title "github.com/satont/twir/apps/parser/internal/commands/channel/title"
	"github.com/satont/twir/apps/parser/internal/commands/clip"
	"github.com/satont/twir/apps/parser/internal/commands/dudes"
	"github.com/satont/twir/apps/parser/internal/commands/games"
	"github.com/satont/twir/apps/parser/internal/commands/manage"
	"github.com/satont/twir/apps/parser/internal/commands/marker"
	"github.com/satont/twir/apps/parser/internal/commands/nuke"
	"github.com/satont/twir/apps/parser/internal/commands/overlays/brb"
	"github.com/satont/twir/apps/parser/internal/commands/overlays/kappagen"
	"github.com/satont/twir/apps/parser/internal/commands/permit"
	"github.com/satont/twir/apps/parser/internal/commands/shoutout"
	"github.com/satont/twir/apps/parser/internal/commands/song"
	sr_youtube "github.com/satont/twir/apps/parser/internal/commands/songrequest/youtube"
	"github.com/satont/twir/apps/parser/internal/commands/spam"
	"github.com/satont/twir/apps/parser/internal/commands/stats"
	"github.com/satont/twir/apps/parser/internal/commands/tts"
	"github.com/satont/twir/apps/parser/internal/types"
)

var defaultCommands = lo.SliceToMap(
	[]*types.DefaultCommand{
		song.CurrentSong,
		channel_game.SetCommand,
		channel_game.History,
		channel_title.SetCommand,
		channel_title.History,
		manage.AddAliaseCommand,
		manage.AddCommand,
		manage.CheckAliasesCommand,
		manage.DelCommand,
		manage.EditCommand,
		manage.RemoveAliaseCommand,
		nuke.Command,
		permit.Command,
		shoutout.ShoutOut,
		spam.Command,
		stats.TopEmotes,
		stats.TopEmotesUsers,
		stats.TopMessages,
		stats.TopPoints,
		stats.TopTime,
		stats.Uptime,
		stats.UserAge,
		stats.UserFollowSince,
		stats.UserFollowage,
		stats.UserMe,
		stats.UserWatchTime,
		tts.DisableCommand,
		tts.EnableCommand,
		tts.PitchCommand,
		tts.RateCommand,
		tts.SayCommand,
		tts.SkipCommand,
		tts.VoiceCommand,
		tts.VoicesCommand,
		tts.VolumeCommand,
		sr_youtube.SkipCommand,
		sr_youtube.SrCommand,
		sr_youtube.SrListCommand,
		sr_youtube.WrongCommand,
		kappagen.Kappagen,
		brb.Start,
		brb.Stop,
		games.EightBall,
		games.RussianRoulette,
		games.Voteban,
		games.Duel,
		games.DuelAccept,
		games.DuelStats,
		games.Seppuku,
		dudes.Jump,
		dudes.Grow,
		dudes.Color,
		dudes.Sprite,
		dudes.Leave,
		seventv.Profile,
		seventv.EmoteFind,
		seventv.EmoteRename,
		seventv.EmoteDelete,
		seventv.EmoteAdd,
		clip.MakeClip,
		marker.Marker,
	}, func(v *types.DefaultCommand) (string, *types.DefaultCommand) {
		return v.Name, v
	},
)

// DefaultCommands is getter for default commands.
func DefaultCommands() map[string]*types.DefaultCommand {
	return defaultCommands
}

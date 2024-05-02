package application

import (
	"context"
	"log"
	"strings"

	"github.com/murasame29/hackathon-util/cmd/config"
	"github.com/murasame29/hackathon-util/pkg/logger"
	"github.com/sourcegraph/conc"
)

type CreateRoleParam struct {
	GuildID       string
	SpreadSheetID string
	Range         string
}

func (as *ApplicationService) CreateRole(ctx context.Context, param CreateRoleParam) error {
	values, err := as.gs.Read(config.Config.Spreadsheets.ID, config.Config.Spreadsheets.Range)
	if err != nil {
		logger.Error(ctx, "failed to read spreadsheet", logger.Field("err", err))
		return err
	}
	var wg conc.WaitGroup
	defer wg.Wait()

	for _, value := range values {
		teamName := value[0]

		if teamName == "欠席" {
			continue
		}

		logger.Info(ctx, "creating role ", logger.Field("role", teamName))

		wg.Go(func() {
			if err := as.ds.CreateRole(ctx, param.GuildID, teamName); err != nil {
				logger.Error(ctx, "failed to create Role", logger.Field("role", teamName))
			}

			logger.Info(ctx, "Role Create Successful", logger.Field("role", teamName))
		})
	}
	return nil
}

type DeleteRoleParam struct {
	GuildID       string
	SpreadSheetID string
	Range         string
}

func (as *ApplicationService) DeleteRole(ctx context.Context, param DeleteRoleParam) error {
	values, err := as.gs.Read(config.Config.Spreadsheets.ID, config.Config.Spreadsheets.Range)
	if err != nil {
		logger.Error(ctx, "Error reading spreadsheet", logger.Field("err", err))
		return err
	}

	roles, err := as.ds.GetRoles(ctx, param.GuildID)
	if err != nil {
		logger.Error(ctx, "Error getting roles", logger.Field("err", err))
		return err
	}

	var wg conc.WaitGroup
	defer wg.Wait()

	for _, value := range values {
		teamName := value[0]
		if teamName == "欠席" {
			continue
		}

		logger.Info(ctx, "delete to role ", logger.Field("role", teamName))

		wg.Go(func() {
			roleID, ok := roles[strings.TrimSpace(teamName)]
			if !ok {
				logger.Error(ctx, "role not found", logger.Field("role", teamName))
				return
			}

			if err := as.ds.DeleteRole(ctx, param.GuildID, roleID); err != nil {
				log.Println(err)
			}

			logger.Info(ctx, "Role Delete Successful", logger.Field("role", teamName))
		})
	}
	return nil
}

type BindRoleParam struct {
	GuildID       string
	SpreadSheetID string
	Range         string
}

func (as *ApplicationService) BindRole(ctx context.Context, param BindRoleParam) error {
	values, err := as.gs.Read(config.Config.Spreadsheets.ID, config.Config.Spreadsheets.Range)
	if err != nil {
		logger.Error(ctx, "failed to read spreadsheet", logger.Field("err", err))
		return err
	}
	users, err := as.ds.GetUsersAll(ctx, config.Config.Discord.GuildID)
	if err != nil {
		logger.Error(ctx, "failed to get users", logger.Field("err", err))
		return err
	}

	roles, err := as.ds.GetRoles(ctx, config.Config.Discord.GuildID)
	if err != nil {
		logger.Error(ctx, "failed to get roles", logger.Field("err", err))
		return err
	}

	var wg conc.WaitGroup
	defer wg.Wait()

	for _, value := range values {
		teamName := value[0]

		// 正規表現に一致する部分を空文字に置き換える
		role := strings.TrimSpace(teamName)

		if role == "欠席" {
			continue
		}

		roleID, ok := roles[role]
		if !ok {
			logger.Error(ctx, "failed to get role", logger.Field("role", role))
			continue
		}

		for _, memer := range value[1:] {
			if memer == "" {
				continue
			}

			userID, ok := users[memer]
			if !ok {
				logger.Error(ctx, "failed to get user", logger.Field("user", memer))
				continue
			}

			wg.Go(func() {
				logger.Info(ctx, "add role", logger.Field("user", memer), logger.Field("role", role))

				if err := as.ds.BindRole(ctx, config.Config.Discord.GuildID, userID, roleID); err != nil {
					logger.Error(ctx, "failed to add role", logger.Field("err", err))
				}

				logger.Info(ctx, "role add successful", logger.Field("user", memer), logger.Field("role", role))
			})
		}
	}

	return nil
}

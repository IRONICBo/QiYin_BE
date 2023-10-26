package utils

import (
	"fmt"

	"github.com/Delta456/box-cli-maker/v2"

	"github.com/IRONICBo/QiYin_BE/internal/config"
)

func Banner() {
	// logo
	logo := `
	______________  ______             __________________
	__  __ \__(_) \/ /__(_)______      ___  __ )__  ____/
	_  / / /_  /__  /__  /__  __ \     __  __  |_  __/   
	/ /_/ /_  / _  / _  / _  / / /     _  /_/ /_  /___   
	\___\_\/_/  /_/  /_/  /_/ /_/      /_____/ /_____/   
																				
APP Mode:
- Version: %s
- Debug: %v
- Log file: %s

Topic: https://www.qiniu.com/activity/detail/651297ed0d50912d3d53307b#topic-introduction
Contributer: @IRONICBo @Baihhh @nnnnn319
`
	content := fmt.Sprintf(logo, config.Config.App.Version, config.Config.App.Debug, config.Config.App.LogFile)

	Box := box.New(box.Config{
		Px:       2,
		Type:     "Round",
		Color:    "Blue",
		TitlePos: "Top",
	})
	Box.Print("QiYin Backend", content)
}

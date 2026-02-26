// About window
package ui

import (
	"fmt"

	"fyne.io/fyne/v2/dialog"
)

func (mw *MainWindow) showAbout() {
	aboutText := `
	Copyright 2026 Matthew Hooper
	This developer, project or application is in no way associated with the device manufacturer.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
	`
	aboutText = fmt.Sprintf("%s Version %s-build.%d\n%s",
		mw.app.Metadata().Name,
		mw.app.Metadata().Version,
		mw.app.Metadata().Build,
		aboutText,
	)

	dialog.ShowInformation("About", aboutText, mw.window)
}

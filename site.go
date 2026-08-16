package campusstore

import "embed"

//go:embed index.php styles.css app.js
var siteFiles embed.FS

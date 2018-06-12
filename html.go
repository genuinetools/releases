package main

const (
	tmpl = `<!DOCTYPE html>
<html lang="en-us">
	<head>
		<meta charset="utf-8">
		<title>GitHub Releases</title>
	</head>
	<body>
		<h1>Latest Releases<h1>

		<table>
			<thead>
				<tr>
					<th>Project</th>
					<th>Release</th>
					<th>Link</th>
				</tr>
			</thead>
			<tbody>
			{{range .}}
				<tr>
					<td>{{.Repository.FullName}}</td>
					<td>{{.Release.TagName}}</td>
					<td>{{.Release.HTMLURL}}</td>
				</tr>
			{{end}}
			</tbody>
		</table>

		<script>
		(function () {
		'use strict';
		var devtools = {
		open: false,
		orientation: null
		};
		var threshold = 160;
		var emitEvent = function (state, orientation) {
		window.dispatchEvent(new CustomEvent('devtoolschange', {
		detail: {
			open: state,
			orientation: orientation
		}
		}));
		};

		setInterval(function () {
		var widthThreshold = window.outerWidth - window.innerWidth > threshold;
		var heightThreshold = window.outerHeight - window.innerHeight > threshold;
		var orientation = widthThreshold ? 'vertical' : 'horizontal';

		if (!(heightThreshold && widthThreshold) &&
				((window.Firebug && window.Firebug.chrome && window.Firebug.chrome.isInitialized) || widthThreshold || heightThreshold)) {
		if (!devtools.open || devtools.orientation !== orientation) {
			emitEvent(true, orientation);
		}

		devtools.open = true;
		devtools.orientation = orientation;
		} else {
		if (devtools.open) {
			emitEvent(false, null);
		}

		devtools.open = false;
		devtools.orientation = null;
		}
		}, 500);

		if (typeof module !== 'undefined' && module.exports) {
		module.exports = devtools;
		} else {
		window.devtools = devtools;
		}
		})();

		</script>
		<script>
		var gt = "                           _               __              __\n" +
				"   ____ ____  ____  __  __(_)___  ___     / /_____  ____  / /____\n" +
				"  / __ ` + "`" + `/ _ \\/ __ \\/ / / / / __ \\/ _ \\   / __/ __ \\/ __ \\/ / ___/\n" +
				" / /_/ /  __/ / / / /_/ / / / / /  __/  / /_/ /_/ / /_/ / (__  )\n" +
				" \\__, /\\___/_/ /_/\\__,_/_/_/ /_/\\___/   \\__/\\____/\\____/_/____/\n" +
				"/____/"

				var printed = false;
		window.addEventListener('devtoolschange', function (e) {
		if(e.detail.open) { printed = true; console.log(gt); }
		});
		</script>

	</body>
</html>`
)

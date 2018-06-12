package main

const (
	tmpl = `<!DOCTYPE html>
<html lang="en-us">
	<head>
		<meta charset="utf-8">
		<title>GitHub Releases</title>
		<style>
			html {
				display: block;
			}
			body{
				font-family: Consolas, Inconsolata, monospace;
				display: block;
				margin: 0;
				font-size: 1rem;
				font-weight: 400;
				line-height: 1.5;
				color: #212529;
				text-align: left;
				background-color: #fff;
			}
			@media (min-width: 1200px) .container {
				max-width: 1140px;
			}
			@media (min-width: 992px) .container {
				max-width: 960px;
			}
			@media (min-width: 768px) .container {
				max-width: 720px;
			}
			@media (min-width: 576px) .container {
				max-width: 540px;
			}
			.container {
				max-width: 100%;
				padding: 1rem;
				margin: auto;
			}
			table {
				background-color: transparent;
				border-color: transparent;
				border-collapse: collapse;
				border-spacing: 2px;
				border-color: grey;
				text-align: inherit;
				font-size: .75rem;
				margin-bottom: 20px;
			}
			thead {
				display: table-header-group;
				vertical-align: middle;
				border-color: inherit;
			}
			tr {
				display: table-row;
				vertical-align: inherit;
				border-color: inherit;
			}
			thead th {
				vertical-align: bottom;
				border-bottom: 2px solid #dee2e6;
			}
			td, th {
				padding: .75rem;
				vertical-align: top;
				border-top: 1px solid #dee2e6;
				display: table-cell;
			}
			tbody {
				display: table-row-group;
				vertical-align: middle;
				border-color: inherit;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h1>Latest Releases</h1>
			<p>This only shows the hashes and download links for linux amd64. For other archs click the tag
			to view the release page.</p>
			<p><small>If you wish to modify this page, the repo is: <a href="https://github.com/genuinetools/releases" target="_blank">genuinetools/releases</a></small></p>

			<table>
				<thead>
					<tr>
						<th>Project</th>
						<th>Release</th>
						<th>download</th>
						<th>sha256</th>
						<th>md5</th>
						<th>download count</th>
					</tr>
				</thead>
				<tbody>
				{{range .}}
					<tr>
						<td><a href="{{.Repository.HTMLURL}}" target="_blank">{{.Repository.FullName}}</a></td>
						<td><a href="{{.Release.HTMLURL}}" target="_blank">{{.Release.TagName}}</a></td>
						<td><a href="{{.BinaryURL}}" target="_blank"><code>{{.BinaryName}}</code></a></td>
						<td><code>{{.BinarySHA256}}</code></td>
						<td><code>{{.BinaryMD5}}</code></td>
						<td><bold>{{.BinaryDownloadCount}}</bold></td>
					</tr>
				{{end}}
				</tbody>
			</table>
		</div>

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

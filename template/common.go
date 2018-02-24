// Code generated by go generate; DO NOT EDIT.
// 2018-02-24 17:47:34.998457627 +0000 GMT

package template

var templateCommonMap = map[string]string{
	"entry_pagination": `{{ define "entry_pagination" }}
<div class="pagination">
    <div class="pagination-prev">
        {{ if .prevEntry }}
            <a href="{{ .prevEntryRoute }}" title="{{ .prevEntry.Title }}" data-page="previous">{{ t "Previous" }}</a>
        {{ else }}
            {{ t "Previous" }}
        {{ end }}
    </div>

    <div class="pagination-next">
        {{ if .nextEntry }}
            <a href="{{ .nextEntryRoute }}" title="{{ .nextEntry.Title }}" data-page="next">{{ t "Next" }}</a>
        {{ else }}
            {{ t "Next" }}
        {{ end }}
    </div>
</div>
{{ end }}`,
	"item_meta": `{{ define "item_meta" }}
<div class="item-meta">
    <ul>
        <li>
            <a href="{{ route "feedEntries" "feedID" .entry.Feed.ID }}" title="{{ .entry.Feed.Title }}">{{ domain .entry.Feed.SiteURL }}</a>
        </li>
        <li>
            <time datetime="{{ isodate .entry.Date }}" title="{{ isodate .entry.Date }}">{{ elapsed .user.Timezone .entry.Date }}</time>
        </li>
        <li>
            <a href="#"
                title="{{ t "Save this article" }}"
                data-save-entry="true"
                data-save-url="{{ route "saveEntry" "entryID" .entry.ID }}"
                data-label-loading="{{ t "Saving..." }}"
                data-label-done="{{ t "Done!" }}"
                >{{ t "Save" }}</a>
        </li>
        <li>
            <a href="{{ .entry.URL }}" target="_blank" rel="noopener noreferrer" referrerpolicy="no-referrer" data-original-link="true">{{ t "Original" }}</a>
        </li>
        <li>
            <a href="#"
                data-toggle-bookmark="true"
                data-bookmark-url="{{ route "toggleBookmark" "entryID" .entry.ID }}"
                data-label-loading="{{ t "Saving..." }}"
                data-label-star="☆ {{ t "Star" }}"
                data-label-unstar="★ {{ t "Unstar" }}"
                data-value="{{ if .entry.Starred }}star{{ else }}unstar{{ end }}"
                >{{ if .entry.Starred }}★ {{ t "Unstar" }}{{ else }}☆ {{ t "Star" }}{{ end }}</a>
        </li>
        <li>
            <a href="#"
                title="{{ t "Change entry status" }}"
                data-toggle-status="true"
                data-label-read="✔ {{ t "Read" }}"
                data-label-unread="✘ {{ t "Unread" }}"
                data-value="{{ if eq .entry.Status "read" }}read{{ else }}unread{{ end }}"
                >{{ if eq .entry.Status "read" }}✘ {{ t "Unread" }}{{ else }}✔ {{ t "Read" }}{{ end }}</a>
        </li>
    </ul>
</div>
{{ end }}`,
	"layout": `{{ define "base" }}
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">

    <meta name="viewport" content="width=device-width, initial-scale=1.0, user-scalable=no">
    <meta name="mobile-web-app-capable" content="yes">
    <meta name="apple-mobile-web-app-title" content="Miniflux">
    <link rel="manifest" href="{{ route "webManifest" }}">

    <meta name="robots" content="noindex,nofollow">
    <meta name="referrer" content="no-referrer">

    <link rel="icon" type="image/png" href="{{ route "appIcon" "filename" "favicon.png" }}">
    <link rel="apple-touch-icon" href="{{ route "appIcon" "filename" "touch-icon-iphone.png" }}">
    <link rel="apple-touch-icon" sizes="72x72" href="{{ route "appIcon" "filename" "touch-icon-ipad.png" }}">
    <link rel="apple-touch-icon" sizes="114x114" href="{{ route "appIcon" "filename" "touch-icon-iphone-retina.png" }}">
    <link rel="apple-touch-icon" sizes="144x144" href="{{ route "appIcon" "filename" "touch-icon-ipad-retina.png" }}">
    <link rel="shortcut icon" type="image/x-icon" href="{{ route "favicon" }}">

    {{ if .csrf }}
        <meta name="X-CSRF-Token" value="{{ .csrf }}">
    {{ end }}
    <title>{{template "title" .}} - Miniflux</title>
    {{ if .user }}
        <link rel="stylesheet" type="text/css" href="{{ route "stylesheet" "name" .user.Theme }}">
    {{ else }}
        <link rel="stylesheet" type="text/css" href="{{ route "stylesheet" "name" "white" }}">
    {{ end }}
    <script type="text/javascript" src="{{ route "javascript" }}" defer></script>
</head>
<body data-entries-status-url="{{ route "updateEntriesStatus" }}">
    {{ if .user }}
    <header class="header">
        <nav>
            <div class="logo">
                <a href="{{ route "unread" }}">Mini<span>flux</span></a>
            </div>
            <ul>
                <li {{ if eq .menu "unread" }}class="active"{{ end }} title="{{ t "Keyboard Shortcut: %s" "g u" }}">
                    <a href="{{ route "unread" }}" data-page="unread">{{ t "Unread" }}</a>
                    {{ if gt .countUnread 0 }}
                        <span class="unread-counter-wrapper">(<span class="unread-counter">{{ .countUnread }}</span>)</span>
                    {{ end }}
                </li>
                <li {{ if eq .menu "starred" }}class="active"{{ end }} title="{{ t "Keyboard Shortcut: %s" "g b" }}">
                    <a href="{{ route "starred" }}" data-page="starred">{{ t "Starred" }}</a>
                </li>
                <li {{ if eq .menu "history" }}class="active"{{ end }} title="{{ t "Keyboard Shortcut: %s" "g h" }}">
                    <a href="{{ route "history" }}" data-page="history">{{ t "History" }}</a>
                </li>
                <li {{ if eq .menu "feeds" }}class="active"{{ end }} title="{{ t "Keyboard Shortcut: %s" "g f" }}">
                    <a href="{{ route "feeds" }}" data-page="feeds">{{ t "Feeds" }}</a>
                </li>
                <li {{ if eq .menu "categories" }}class="active"{{ end }} title="{{ t "Keyboard Shortcut: %s" "g c" }}">
                    <a href="{{ route "categories" }}" data-page="categories">{{ t "Categories" }}</a>
                </li>
                <li {{ if eq .menu "settings" }}class="active"{{ end }} title="{{ t "Keyboard Shortcut: %s" "g s" }}">
                    <a href="{{ route "settings" }}" data-page="settings">{{ t "Settings" }}</a>
                </li>
                <li>
                    <a href="{{ route "logout" }}" title="{{ t "Logged as %s" .user.Username }}">{{ t "Logout" }}</a>
                </li>
            </ul>
        </nav>
    </header>
    {{ end }}
    {{ if .flashMessage }}
        <div class="flash-message alert alert-success">{{ .flashMessage }}</div>
    {{ end }}
    {{ if .flashErrorMessage }}
        <div class="flash-error-message alert alert-error">{{ .flashErrorMessage }}</div>
    {{ end }}
    <main>
        {{template "content" .}}
    </main>
    <template id="keyboard-shortcuts">
        <div id="modal-left">
            <a href="#" class="btn-close-modal">x</a>
            <h3>{{ t "Keyboard Shortcuts" }}</h3>

            <div class="keyboard-shortcuts">
                <p>{{ t "Sections Navigation" }}</p>
                <ul>
                    <li>{{ t "Go to unread" }} = <strong>g + u</strong></li>
                    <li>{{ t "Go to bookmarks" }} = <strong>g + b</strong></li>
                    <li>{{ t "Go to history" }} = <strong>g + h</strong></li>
                    <li>{{ t "Go to feeds" }} = <strong>g + f</strong></li>
                    <li>{{ t "Go to categories" }} = <strong>g + c</strong></li>
                    <li>{{ t "Go to settings" }} = <strong>g + s</strong></li>
                    <li>{{ t "Show keyboard shortcuts" }} = <strong>?</strong></li>
                </ul>

                <p>{{ t "Items Navigation" }}</p>
                <ul>
                    <li>{{ t "Go to previous item" }} = <strong>p or j or ◄</strong></li>
                    <li>{{ t "Go to next item" }} = <strong>n or k or ►</strong></li>
                </ul>

                <p>{{ t "Pages Navigation" }}</p>
                <ul>
                    <li>{{ t "Go to previous page" }} = <strong>h</strong></li>
                    <li>{{ t "Go to next page" }} = <strong>l</strong></li>
                </ul>

                <p>{{ t "Actions" }}</p>
                <ul>
                    <li>{{ t "Open selected item" }} = <strong>o</strong></li>
                    <li>{{ t "Open original link" }} = <strong>v</strong></li>
                    <li>{{ t "Toggle read/unread" }} = <strong>m</strong></li>
                    <li>{{ t "Mark current page as read" }} = <strong>A</strong></li>
                    <li>{{ t "Download original content" }} = <strong>d</strong></li>
                    <li>{{ t "Toggle bookmark" }} = <strong>f</strong></li>
                    <li>{{ t "Save article" }} = <strong>s</strong></li>
                    <li>{{ t "Close modal dialog" }} = <strong>Esc</strong></li>
                </ul>
            </div>
        </div>
    </template>
</body>
</html>
{{ end }}
`,
	"pagination": `{{ define "pagination" }}
<div class="pagination">
    <div class="pagination-prev">
        {{ if .ShowPrev }}
            <a href="{{ .Route }}{{ if gt .PrevOffset 0 }}?offset={{ .PrevOffset }}{{ end }}" data-page="previous">{{ t "Previous" }}</a>
        {{ else }}
            {{ t "Previous" }}
        {{ end }}
    </div>

    <div class="pagination-next">
        {{ if .ShowNext }}
            <a href="{{ .Route }}?offset={{ .NextOffset }}" data-page="next">{{ t "Next" }}</a>
        {{ else }}
            {{ t "Next" }}
        {{ end }}
    </div>
</div>
{{ end }}
`,
}

var templateCommonMapChecksums = map[string]string{
	"entry_pagination": "f1465fa70f585ae8043b200ec9de5bf437ffbb0c19fb7aefc015c3555614ee27",
	"item_meta":        "4796b74adca0567f3dbf8bdf6ac8cda59f455ea34cb6d4a92c83660fa72a883d",
	"layout":           "c7565e2cf904612e236bc1d7167c6c124ffe5d27348608eb5c2336606f266896",
	"pagination":       "6ff462c2b2a53bc5448b651da017f40a39f1d4f16cef4b2f09784f0797286924",
}

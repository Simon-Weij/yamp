// qmllint disable unqualified
pragma ComponentBehavior: Bound
import QtQuick
import QtQuick.Controls
import ".."

TextField {
    id: searchField
    property string searchMode: "Songs"
    signal resultsReady(var results)
    width: 500
    leftPadding: 10
    topPadding: 10
    bottomPadding: 10
    placeholderTextColor: Theme.foregroundSubtle
    placeholderText: "Search..."
    background: Rectangle {
        color: Theme.backgroundSecondary
        radius: 12
    }
    color: Theme.foreground
    onTextChanged: debounceTimer.restart()
    onSearchModeChanged: {
        if (searchField.text == "" || searchField.text == null || searchField.text == undefined) {
            searchField.resultsReady([])
            return;
        }

        debounceTimer.restart()
    }

    Timer {
        id: debounceTimer
        interval: 750
        repeat: false
        onTriggered: {
            if (searchField.text == "" || searchField.text == null || searchField.text == undefined) {
                searchField.resultsReady([])
                return;
            }

            const results = searchField.searchMode === "Albums"
                ? api.searchAlbums(searchField.text)
                : api.searchSongs(searchField.text);
            searchField.resultsReady(results)
        }
    }
}

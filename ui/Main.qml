// qmllint disable unqualified
pragma ComponentBehavior: Bound
import QtQuick
import QtQuick.Controls
import "."
import "components"

ApplicationWindow {
    visible: true
    width: 500
    height: 500
    title: "yamp"
    color: Theme.backgroundPrimary

    property var searchResults: []
    property string currentFilter: "Songs"

    SearchField {
        id: searchField
        width: parent.width - 40
        anchors.horizontalCenter: parent.horizontalCenter
        anchors.top: parent.top
        anchors.margins: 20
        searchMode: currentFilter
        onResultsReady: function (results) {
            searchResults = results;
        }
    }

    FilterBar {
        id: filterBar
        anchors.top: searchField.bottom
        anchors.horizontalCenter: searchField.horizontalCenter
        anchors.topMargin: 16
        onFilterChanged: function (filter) {
            currentFilter = filter;
        }
    }

    ListView {
        anchors.left: parent.left
        anchors.right: parent.right
        anchors.top: filterBar.bottom
        anchors.bottom: parent.bottom
        anchors.topMargin: 16
        anchors.leftMargin: 20
        anchors.rightMargin: 20
        anchors.bottomMargin: 20
        clip: true
        spacing: 12
        model: searchResults

        delegate: ResultCard {
            required property var modelData
            itemData: modelData
        }
    }
}

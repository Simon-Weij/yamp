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

        delegate: Rectangle {
            required property var modelData
            width: ListView.view.width
            height: 120
            radius: 16
            color: Theme.backgroundSecondary
            border.color: Theme.foregroundSubtle
            border.width: 1

            Row {
                anchors.fill: parent
                anchors.margins: 16
                spacing: 16

                Rectangle {
                    width: 88
                    height: 88
                    radius: 12
                    clip: true
                    color: Theme.backgroundPrimary

                    Rectangle {
                        anchors.fill: parent
                        radius: parent.radius
                        color: Theme.backgroundPrimary
                        visible: coverArt.status !== Image.Ready

                        SequentialAnimation on opacity {
                            loops: Animation.Infinite
                            NumberAnimation {
                                from: 0.45
                                to: 0.9
                                duration: 700
                                easing.type: Easing.InOutQuad
                            }
                            NumberAnimation {
                                from: 0.9
                                to: 0.45
                                duration: 700
                                easing.type: Easing.InOutQuad
                            }
                        }
                    }

                    Rectangle {
                        anchors.fill: parent
                        radius: parent.radius
                        clip: true
                        color: "transparent"

                        Image {
                            id: coverArt
                        anchors.fill: parent
                        source: modelData.release_id
                            ? "https://coverartarchive.org/release/" + modelData.release_id + "/front"
                            : ""
                        asynchronous: true
                        cache: true
                        fillMode: Image.PreserveAspectCrop
                        visible: status === Image.Ready
                        }
                    }
                }

                Column {
                    anchors.verticalCenter: parent.verticalCenter
                    width: parent.width - 104
                    spacing: 8

                    Text {
                        width: parent.width
                        text: modelData.title || "Unknown result"
                        color: Theme.foreground
                        font.pixelSize: 20
                        font.bold: true
                        elide: Text.ElideRight
                    }

                    Text {
                        width: parent.width
                        text: modelData.artist || "Unknown artist"
                        color: Theme.foregroundSubtle
                        font.pixelSize: 14
                        elide: Text.ElideRight
                    }
                }
            }
        }
    }
}

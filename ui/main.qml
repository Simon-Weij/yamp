// qmllint disable unqualified

pragma ComponentBehavior: Bound
import QtQuick
import QtQuick.Controls
import "."

ApplicationWindow {
    visible: true
    width: 500
    height: 500
    title: "yamp"
    color: Theme.backgroundPrimary

    TextField {
        id: searchField
        anchors.horizontalCenter: parent.horizontalCenter
        anchors.top: parent.top
        anchors.margins: 20
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

        Timer {
            id: debounceTimer
            interval: 750
            repeat: false
            onTriggered: {
                if (searchField.text == "" || searchField.text == null || searchField.text == undefined)
                    return;

                const results = api.searchRecordings(searchField.text);
                for (const r of results) {
                    console.log("id: " + r.id);
                    console.log("title: " + r.title);
                    console.log("artist: " + r.artist);
                    console.log("album: " + r.album);
                    console.log("duration: " + r.duration_ms);
                    console.log("release date: " + r.release_date);
                    console.log("release id: " + r.release_id);
                    console.log("release group id: " + r.release_group_id);
                    console.log("-".repeat(60));
                }
            }
        }
    }
    Row {
        anchors.top: searchField.bottom
        anchors.horizontalCenter: searchField.horizontalCenter
        anchors.topMargin: 16
        spacing: 16

        ButtonGroup {
            id: filterGroup
        }

        Repeater {
            model: ["Songs", "Albums", "Artists"]
            Button {
                id: btn
                required property string modelData
                text: modelData
                checkable: true
                leftPadding: 48
                topPadding: 16
                bottomPadding: 16
                rightPadding: 48
                checked: modelData === "Songs"
                ButtonGroup.group: filterGroup

                HoverHandler {
                    cursorShape: Qt.PointingHandCursor
                }

                background: Rectangle {
                    radius: 6
                    color: btn.checked ? Theme.foreground : Theme.backgroundSecondary
                }
                contentItem: Text {
                    text: btn.text
                    color: btn.checked ? Theme.backgroundPrimary : Theme.foreground
                    horizontalAlignment: Text.AlignHCenter
                    verticalAlignment: Text.AlignVCenter
                }
            }
        }
    }
}

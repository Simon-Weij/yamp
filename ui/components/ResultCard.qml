// qmllint disable unqualified
pragma ComponentBehavior: Bound
import QtQuick
import ".."

Rectangle {
    required property var itemData

    function formatDuration(durationMs) {
        const totalSeconds = Math.floor((durationMs || 0) / 1000);
        const hours = Math.floor(totalSeconds / 3600);
        const minutes = Math.floor((totalSeconds % 3600) / 60);
        const seconds = totalSeconds % 60;

        if (hours > 0)
            return hours + ":" + String(minutes).padStart(2, "0") + ":" + String(seconds).padStart(2, "0");

        return minutes + ":" + String(seconds).padStart(2, "0");
    }

    function playSong(data) {
        api.playSong(data);
    }

    width: ListView.view.width
    height: 120
    radius: 16
    color: Theme.backgroundSecondary
    border.color: Theme.foregroundSubtle
    border.width: 1

    MouseArea {
        anchors.fill: parent
        acceptedButtons: Qt.LeftButton
        onDoubleClicked: playSong(itemData)
    }

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
                    source: itemData.release_id ? "https://coverartarchive.org/release/" + itemData.release_id + "/front" : ""
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
                text: itemData.title || "Unknown result"
                color: Theme.foreground
                font.pixelSize: 20
                font.bold: true
                elide: Text.ElideRight
            }

            Text {
                width: parent.width
                text: itemData.album || "Unknown album"
                color: Theme.foregroundSubtle
                font.pixelSize: 14
                elide: Text.ElideRight
            }

            Text {
                width: parent.width
                text: itemData.artist || "Unknown artist"
                color: Theme.foregroundSubtle
                font.pixelSize: 14
                elide: Text.ElideRight
            }

            Text {
                width: parent.width
                text: formatDuration(itemData.duration_ms)
                color: Theme.foregroundSubtle
                font.pixelSize: 13
                visible: (itemData.duration_ms || 0) > 0
            }
        }
    }
}

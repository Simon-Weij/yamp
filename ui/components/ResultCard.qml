// qmllint disable unqualified
pragma ComponentBehavior: Bound
import QtQuick
import QtQuick.Controls
import ".."

Rectangle {
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
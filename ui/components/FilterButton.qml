// qmllint disable unqualified
pragma ComponentBehavior: Bound
import QtQuick
import QtQuick.Controls
import ".."

Button {
    id: btn
    checkable: true
    leftPadding: 48
    topPadding: 16
    bottomPadding: 16
    rightPadding: 48

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

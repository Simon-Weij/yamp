import QtQuick
import QtQuick.Controls
import "."

ApplicationWindow {
    visible: true
    width: 500
    height: 500
    title: "yamp"
    color: Theme.background

    Text {
        text: "Hello world"
        color: Theme.foreground
        font.pixelSize: Theme.fontSizeMedium
    }
}

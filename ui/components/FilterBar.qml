// qmllint disable unqualified
pragma ComponentBehavior: Bound
import QtQuick
import QtQuick.Controls

Row {
    signal filterChanged(string filter)
    spacing: 16

    ButtonGroup {
        id: filterGroup
    }

    Repeater {
        model: ["Songs", "Albums"]
        FilterButton {
            required property string modelData
            text: modelData
            checked: modelData === "Songs"
            ButtonGroup.group: filterGroup
            onClicked: filterBar.filterChanged(modelData)
        }
    }
}

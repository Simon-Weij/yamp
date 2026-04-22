import sys

from PySide6.QtCore import QFileSystemWatcher, QUrl
from PySide6.QtGui import QGuiApplication
from PySide6.QtQml import QQmlApplicationEngine


def main():
    app = QGuiApplication(sys.argv)

    engine = QQmlApplicationEngine()
    qml = "./ui/main.qml"

    def load():
        for root in engine.rootObjects():
            root.deleteLater()

        engine.clearComponentCache()
        engine.load(QUrl.fromLocalFile(qml))

    load()

    watcher = QFileSystemWatcher([qml])
    watcher.fileChanged.connect(lambda _: (load(), watcher.addPath(qml)))

    if not engine.rootObjects():
        sys.exit(-1)

    sys.exit(app.exec())


if __name__ == "__main__":
    main()

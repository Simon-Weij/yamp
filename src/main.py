import os
import sys

from PySide6.QtCore import QFileSystemWatcher, QUrl
from PySide6.QtGui import QGuiApplication
from PySide6.QtQml import QQmlApplicationEngine

app = QGuiApplication(sys.argv)
engine = QQmlApplicationEngine()
engine.addImportPath("./ui")


def load():
    engine.clearComponentCache()
    [r.deleteLater() for r in engine.rootObjects()]
    engine.load(QUrl.fromLocalFile("./ui/main.qml"))


watcher = QFileSystemWatcher()
watcher.fileChanged.connect(lambda p: (watcher.addPath(p), load()))


def watch():
    for root, _, files in os.walk("./ui"):
        watcher.addPaths(
            [f"{root}/{f}" for f in files if f.endswith((".qml", "qmldir"))]
        )


watch()
load()
sys.exit(app.exec())

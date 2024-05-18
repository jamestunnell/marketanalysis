import van from "vanjs-core"

const {i} = van.tags

const IconAdd = () => {
    return i({class: "fa-solid fa-plus"});
}

const IconCheck = () => {
    return i({class: "fa-solid fa-check"});
}

const IconError = () => {
    return i({class: "fa-solid fa-triangle-exclamation"});
}

const IconDelete = () => {
    return i({class: "fa-solid fa-trash"});
}

const IconExport = () => {
    return i({class: "fa-solid fa-file-export"});
}

const IconImport = () => {
    return i({class: "fa-solid fa-file-import"});
}

const IconPlay = () => {
    return i({class: "fa-solid fa-play"});
}

const IconSave = () => {
    return i({class: "fa-solid fa-floppy-disk"});
}

const IconView = () => {
    return i({class: "fa-solid fa-eye"});
}

export {IconAdd, IconCheck, IconDelete, IconError, IconExport, IconImport, IconPlay, IconSave, IconView};
import van from "vanjs-core"

const {i} = van.tags

const IconAdd = () => {
    return i({class: "fa-solid fa-plus"});
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

const IconSave = () => {
    return i({class: "fa-solid fa-floppy-disk"});
}

const IconView = () => {
    return i({class: "fa-solid fa-eye"});
}

export {IconAdd, IconDelete, IconExport, IconImport, IconSave, IconView};
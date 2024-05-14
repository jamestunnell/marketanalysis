import van from "vanjs-core"

const { button } = van.tags

const ButtonAdd = (onclick) => {
    const btn = ButtonAct({text: "", onclick: onclick});

    btn.classList.add("fa-solid", "fa-plus");
}

const ButtonDelete = (onclick) => {
    const btn = ButtonAct({text: "", onclick: onclick});

    btn.classList.add("fa-solid", "fa-trash");

    return btn
}

const ButtonExport = (onclick) => {
    const btn = ButtonAct({text: "", onclick: onclick});

    btn.classList.add("fa-solid", "fa-file-export");
}

const ButtonImport = (onclick) => {
    const btn = ButtonAct({text: "", onclick: onclick});

    btn.classList.add("fa-solid", "fa-file-import");
}

const ButtonSave = (onclick) => {
    const btn = ButtonAct({text: "", onclick: onclick});

    btn.classList.add("fa-solid", "fa-floppy-disk");

    return btn
}

const ButtonView = (onclick) => {
    const btn = ButtonAct({text: "", onclick: onclick});

    btn.classList.add("fa-regular", "fa-eye");
    
    return btn
}

const ButtonAct = ({text, onclick}) => {
    return button(
        {
            class: `bg-indigo-500 p-3 m-1 rounded-md text-white hover:bg-indigo-400`,
            onclick: onclick,
        },
        text,
    );
}

const ButtonCancel = ({text, onclick}) => {
    return button(
        {
            class: "bg-gray-100 p-3 m-1 rounded-md text-black hover:bg-gray-200",
            onclick: onclick,
        },
        text,
    )
}

export { ButtonAdd, ButtonDelete, ButtonExport, ButtonImport, ButtonSave, ButtonView, ButtonAct, ButtonCancel };
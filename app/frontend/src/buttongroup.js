import van from "vanjs-core"

const {div} = van.tags

const ButtonGroup = ({buttons, moreClass=""}) => {
    return buttonGroup({buttons, moreClass, hidden: false})
}

const ButtonGroupHideable = ({buttons, hidden, moreClass=""}) => {
    return buttonGroup({buttons, moreClass, hidden})
}

const buttonGroup = ({buttons, moreClass="", hidden}) => {
    let classStr = "flex flex-row space-x-2"
    if (moreClass.length !== 0) {
        classStr += ` ${moreClass}`
    }

    return div({class: classStr, hidden}, buttons);
}

export { ButtonGroup, ButtonGroupHideable }
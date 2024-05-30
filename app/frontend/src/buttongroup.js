import van from "vanjs-core"

const {div} = van.tags

const ButtonGroup = ({buttons, moreClass=""}) => {
    return buttonGroup({buttons, moreClass, hidden: false})
}

const ButtonGroupHideable = ({buttons, hidden, moreClass=""}) => {
    const group = buttonGroup({buttons, moreClass})

    van.derive(() => {
        if (hidden.val) {
            group.classList.add("hidden")
        } else {
            group.classList.remove("hidden")
        }
    })
    
    return group
}

const buttonGroup = ({buttons, moreClass="",}) => {
    let classStr = "flex flex-row space-x-2"
    if (moreClass.length !== 0) {
        classStr += ` ${moreClass}`
    }

    return div({class: classStr}, buttons);
}

export { ButtonGroup, ButtonGroupHideable }
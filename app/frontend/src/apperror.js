import van from "vanjs-core"

const { div, h2, p } = van.tags

const AppError = ({type, msg, details, hidden}) => {
    const divClass = van.derive(() => {
        return `text-red-500 space-y-6 ${hidden.val ? "hidden" : ""}`
    });
    const title = van.derive(() => {
        switch (type.val) {
            case "InvalidInput":
                return "Invalid Input Error"
            case "NotFound":
                return "Not Found Error"
            case "ActionFailed":
                return "Server Action Error"
        }

        return "Unknown Error"
    });
    const detailsCombined = van.derive(() => details.val.join(", "));

    return div(
        {class: divClass},
        h2({class: "text-2xl font-bold"}, title),
        div(
            {class: "grid grid-cols-2 gap-4"},
            "Message",
            p(msg),
            "Details",
            p(detailsCombined),
        ),
    )
};

export default AppError;
import van from "vanjs-core"

const {a, button, div, h3, input, li, option, select, ul} = van.tags

const BASE_URL = `https://4002-debug-jamestunnel-marketanaly-7v91pin8jv5.ws-us110.gitpod.io`
const TIME_ZONES = ["America/New_York", "America/Los_Angeles", "Greenwich"]

const getSecurities = async () => {
    console.log("getting securities");

    const resp = await fetch(`${BASE_URL}/securities`, {credentials: 'same-origin'});

    if (!resp.ok) {
        console.log("non-ok get securities response %d: %s", resp.status, resp.text());

        return []
    }
    
    const d = await resp.json()

    console.log("get securities response JSON: %o", d)

    return d.securities;
}

const addSecurity = async (item) => {
    console.log("adding security");

    const resp = await fetch(`${BASE_URL}/securities`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json; charset=utf-8'
        },
        body: JSON.stringify(item),
        credentials: 'same-origin'
    });

    console.log('status:', response.status)
}

const AddForm = () => {
    const sym = van.state("");
    const tz = van.state("America/New_York");
    const open = van.state("09:30");
    const close = van.state("16:00");

    return div(
        "Symbol: ",
        input({type: "text", value: sym, oninput: e => sym.val = e.target.value}),
        "Time Zone :",
        select({oninput: e => tz = e.target.value, value: tz},
            TIME_ZONES.map(x => option({value: x}, x))),
        input({type: "text", value: open, oninput: e => open.val = e.target.value}),
        input({type: "text", value: close, oninput: e => close.val = e.target.value}),
        button({onclick: () => {
            addSecurity({symbol: sym.val, timeZone: tz.val, open: open.val, close: close.val})
        }}, "➕"),
    )
}

const ListItem = ({text}) => {
    const deleted = van.state(false)
    return () => deleted.val ? null : li(
        div(
            {class: "flex flex-row gap-4"},
            text, a({onclick: () => deleted.val = true}, "❌"),
        )
    )
}

const Securities = () => {
    const listDom = ul({class:"p-4"})

    getSecurities().then(
        (items) => {
            console.log("found %d securities", items.length);
            
            const listItems = items.map(x => ListItem({text: x.symbol}))
            
            van.add(listDom, listItems)
        }
    ).catch(error => {
        console.log("failed to get securities: " + error);
    });

    const dom = div();

    return div(
        {class: "h-screen w-screen p-4"},
        dom,
        h3("Securities"),
        AddForm(),
        listDom
    );
}

van.add(document.body, Securities())
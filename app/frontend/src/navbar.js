import van from "vanjs-core"
import { routeTo } from 'vanjs-router'

const {  button, div, li, nav, ul } = van.tags

const NavBarItem = ({text, route, routeArgs, currentRoute}) => {
    const textCls = currentRoute.includes(route) ? "bg-gray-400": "hover:bg-gray-600";
    
    return li(
        {class: `md:px-4 md:py-2 ${textCls}`},
        button({onclick: () => { routeTo(route, routeArgs) }}, text),
    )
}

const NavBar = ({currentRoute}) => {
    return div(
        nav(
            {class: "nav bg-gray-500 text-white"},
            div(
                {class: "container flex order-3 w-full"},
                ul(
                    {class: "flex flex-wrap items-center text-lg font-semibold"},
                    NavBarItem({text: 'Home', route: 'home', routeArgs: [], currentRoute}),
                    NavBarItem({text: 'Graphs', route: 'graphs', routeArgs: [], currentRoute}),
                )
            )
        )
    )
}

export default NavBar;

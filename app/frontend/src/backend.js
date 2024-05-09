const BASE_URL = import.meta.env.VITE_API_BASE_URL;

const Delete = async (route) => {
    return fetch(BASE_URL + route, {method: "DELETE"});
}

const Get = async (route) => {
    return fetch(
        BASE_URL + route,
        {
            method: "GET",
            headers: {
                'Accept': 'application/json',
            }
        },
    );
}

const Post = async ({route, content}) => {
    return fetch(
        BASE_URL + route,
        {
            method: "POST",
            headers: {
                'Content-Type': 'application/json; charset=utf-8'
            },
            body: JSON.stringify(content),
        },
    );
}

const Put = async ({route, content}) => {
    return fetch(
        BASE_URL + route,
        {
            method: "PUT",
            headers: {
                'Content-Type': 'application/json; charset=utf-8'
            },
            body: JSON.stringify(content),
        },
    );
}

export {Delete, Get, Post, Put};
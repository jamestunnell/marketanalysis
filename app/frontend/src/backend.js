const BASE_URL = '/api';

const Delete = async (route) => {
    const url = BASE_URL + route;

    console.log("fetch DELETE %s", url);

    return fetch(url, {method: "DELETE"});
}

const Get = async (route) => {
    const url = BASE_URL + route;

    console.log("fetch GET %s", url);

    return fetch(
        url,
        {
            method: "GET",
            headers: {
                'Accept': 'application/json',
            }
        },
    );
}

const Post = async ({route, content}) => {
    const url = BASE_URL + route;

    console.log("fetch POST %s", url);

    return fetch(
        url,
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
    const url = BASE_URL + route;

    console.log("fetch PUT %s", url);

    return fetch(
        url,
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
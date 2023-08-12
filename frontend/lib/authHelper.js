import axios from "axios"

const AxiosInstance = axios.create({
    baseURL: process.env.NEXT_PUBLIC_APP_HOST,
    timeout: 10000,
})

export async function GetCsrf() {
    const response = await axios.get(process.env.NEXT_PUBLIC_APP_HOST + "/api/auth/csrf")
    return response.headers.get("X-Csrf-Token")
}

AxiosInstance.interceptors.request.use(
    async (config) => {
        let csrfToken = JSON.parse(localStorage.getItem("csrfToken"))
        if (!csrfToken) {
            csrfToken = await GetCsrf()
            localStorage.setItem("X-Csrf-Token", csrfToken)
        }

        config.headers.set("X-Csrf-Token", csrfToken)
        return config
    }
)

export async function CreateSession(idToken) {
    return AxiosInstance.post(process.env.NEXT_PUBLIC_APP_HOST + "/api/auth/sign-in", {
        id_token: idToken,
    });
}

export async function EndSession() {
    return AxiosInstance.post(process.env.NEXT_PUBLIC_APP_HOST + "/api/auth/sign-out", {
        withCredentials: true,
    });
}

export default AxiosInstance;
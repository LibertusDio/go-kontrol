import axios from "axios";

const apiClient = () => {
    const {REACT_APP_API_URL, ACCESS_TOKEN_STORAGE} = process.env;
    const axiosInstance = axios.create({
        baseURL: REACT_APP_API_URL,
        responseType: "json",
        validateStatus: (status) => status < 500
    });
    const accessToken = localStorage.getItem(
        ACCESS_TOKEN_STORAGE ? ACCESS_TOKEN_STORAGE : "access_token"
    );
    if (accessToken) {
        axiosInstance.interceptors.request.use(function (config) {
            config.headers!.Authorization = accessToken
                ? `Bearer ${accessToken}`
                : "";
            return config;
        });
    }
    axiosInstance.interceptors.response.use(
        (response) => {
            if (response.status !== 200) {
                response.data = Object.assign(response.data, {code: response.status, message: response.statusText})
            }
            return response
        },
        (error) => {
            return error
        }
    );
    return axiosInstance;
};

export default apiClient;

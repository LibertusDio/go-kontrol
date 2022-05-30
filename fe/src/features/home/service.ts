import axios, {AxiosError, AxiosResponse} from "axios";
import API_URL from '../../constants';

axios.interceptors.response.use(function (response: AxiosResponse) {
    return response;
}, function (error:AxiosError) {
    // @ts-ignore
    return Promise.reject({code: error.response?.status, message: error.response?.statusText});
});

const serviceInfo = (serviceId: string) => {
    const token = localStorage.getItem('token')
    let headers = {}
    if (token) {
        headers = {
            Authorization:  token ? 'Bearer ' + token : '',
        }
    }
    return axios
        .get(API_URL + `/${serviceId}/api/info`, {
            headers: headers,
        })
}
const services = {
    serviceInfo
}
export default services

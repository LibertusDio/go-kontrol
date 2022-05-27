import axios from "axios";
import API_URL from '../../constants';

const serviceInfo = (serviceId: string) => {
    return axios
        .get(API_URL + `/${serviceId}/api/info`, {
            headers: {
                Authorization: 'Bearer ' + localStorage.getItem('token'),
                // 'Access-Control-Allow-Origin': '*',
                // 'Content-Type': 'application/json;charset=utf-8',
            },
        });
}
const services = {
    serviceInfo
}
export default services

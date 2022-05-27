import {createAsyncThunk} from "@reduxjs/toolkit";
import services from "./service";


export const serviceInfo = createAsyncThunk("service/info", async (serviceId: string) => {
    return await services.serviceInfo(serviceId)
})

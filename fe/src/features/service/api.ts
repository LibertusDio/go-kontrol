import {createAsyncThunk} from "@reduxjs/toolkit";
import apiClient from "../../helper/apiClient";


export const getServiceInfo = createAsyncThunk("get/service/info", async (id: string, { rejectWithValue }) => {
    try {
        const resp = await apiClient().get(`/${id}/api/info`)
        if (resp.status !== 200) {
            return rejectWithValue(resp.data)
        }
        return resp.data;
    }
    catch (error) {
        return rejectWithValue(error);
    }
});

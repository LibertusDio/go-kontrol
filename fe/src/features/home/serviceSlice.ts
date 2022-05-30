import {createAsyncThunk, createSlice} from "@reduxjs/toolkit";
import services from "./service";
import {RootState} from "../../app/store";


export const serviceInfo = createAsyncThunk("service/info", async (serviceId: string, thunkAPI: any) => {
    try {
        const res = await services.serviceInfo(serviceId)
        if (res.status !== 200) {
            return thunkAPI.rejectWithValue(res.data)
        }
        return res.data
    } catch (e) {
        return thunkAPI.rejectWithValue(e)
    }

})

export interface ResponseServiceState {
    statusCode: number
    message: string
    effectedAt: number
}

const initialServiceState: ResponseServiceState = {
    statusCode: 0,
    message: '',
    effectedAt: (new Date()).getTime(),
}

const serviceSlice = createSlice({
    name: "service/info",
    extraReducers: {
        // @ts-ignore
        [serviceInfo.fulfilled]: (state: ResponseServiceState, {payload}) => {
            const p = payload as { code: number, message: string }
            state.statusCode = p.code
            state.message = p.message
            state.effectedAt = (new Date()).getTime()
        },
        // @ts-ignore
        [serviceInfo.rejected]: (state: ResponseServiceState, {payload}) => {
            const p = payload as { code: number, message: string }
            state.statusCode = p.code
            state.message = p.message
            state.effectedAt = (new Date()).getTime()
        },
        // @ts-ignore
        [serviceInfo.pending]: (state: ResponseServiceState, {payload}) => {
            console.log("pending: ", state, payload)
        },
    },
    initialState: initialServiceState,
    reducers: {}

})

export default serviceSlice.reducer;
export const serviceSelector = (state: RootState) => state.serviceInfo;

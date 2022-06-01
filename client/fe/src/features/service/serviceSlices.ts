import {createSlice, PayloadAction} from "@reduxjs/toolkit";
import {getServiceInfo} from "./api";
import {RootState} from "../../app/store";

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
    extraReducers: {
        // @ts-ignore
        [getServiceInfo.fulfilled]: (state: ResponseServiceState, action: PayloadAction<any>) => {
            state.statusCode = action.payload?.code
            state.message = action.payload?.message
        },
        // @ts-ignore
        [getServiceInfo.rejected]: (state: ResponseServiceState, action: PayloadAction<any>) => {
            state.statusCode = action.payload?.code
            state.message = action.payload?.message
        },
    },
    initialState: initialServiceState,
    name: "service",
    reducers: {}

})

export default serviceSlice.reducer

export const selectService = (state: RootState) => state.service

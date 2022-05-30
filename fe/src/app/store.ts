import {configureStore, ThunkAction, Action} from '@reduxjs/toolkit';
import counterReducer from '../features/counter/counterSlice';
import authReducer from '../features/login/login-slice'
import serviceReducer from "../features/home/serviceSlice";

export const store = configureStore({
    reducer: {
        counter: counterReducer,
        auth: authReducer,
        serviceInfo: serviceReducer
    },
});

export type AppDispatch = typeof store.dispatch;
export type RootState = ReturnType<typeof store.getState>;
export type AppThunk<ReturnType = void> = ThunkAction<ReturnType,
    RootState,
    unknown,
    Action<string>>;

export interface Permission {
    object_id: string
    token: string
}

export interface AuthModel {
    code: number,
    message: string,
    object_permission: Permission
}

export interface AuthRequestInterface {
    service_id: string
    user_name: string
    password: string
}

export type AuthRequest = {
    service_id: string
    user_name: string
    password: string
}

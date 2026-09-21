import { http } from "@/utils/request";

export const uploadApi = {
  uploadImage: (file: File) => {
    const form = new FormData();
    form.append("file", file);
    return http.upload<{ url: string }>("/uploads?kind=image", form);
  },
  uploadAudio: (file: File) => {
    const form = new FormData();
    form.append("file", file);
    return http.upload<{ url: string }>("/uploads?kind=audio", form);
  },
};

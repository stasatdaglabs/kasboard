import dateFormat from "dateformat";

export const timestampToString = (timestamp: number) => {
    const date = new Date(timestamp);
    return dateFormat(date, "yyyy-mm-dd HH:MM:ss");
};

import re
from datetime import datetime


def validate_month(month):
    return re.match(r"[0-9]{4}-([0][0-9]|1[0-2])", month)


def validate_date(date):
    return re.match(r"[0-9]{4}-([0][0-9]|1[0-2])-([0-2][0-9]|3[01])", date)


def get_first_monday(month):
    datetime_start = datetime.strptime(month + "-01", "%Y-%m-%d")
    datetime_monday = datetime.fromtimestamp(datetime_start.timestamp() - datetime_start.weekday()*86400)
    return datetime_monday.strftime("%Y-%m-%d")


def get_last_monday(month):
    end = get_last_day_in_month(month)
    datetime_end = datetime.strptime(end, "%Y-%m-%d")
    datetime_friday = datetime.fromtimestamp(datetime_end.timestamp() + (7 - datetime_end.weekday())*86400)
    return datetime_friday.strftime("%Y-%m-%d")


def get_last_day_in_month(month):
    month_only = month[5:7]
    if month_only in ["01", "03", "05", "07", "08", "10", "12"]:
        end = month + "-31"
    elif month_only == "02":
        if leap_year(month[0:4]):
            end = month + "-29"
        else:
            end = month + "-28"
    else:
        end = month + "-30"
    return end


def leap_year(year):
    if (year % 4) == 0:
        if (year % 100) == 0:
            if (year % 400) == 0:
                return True
            else:
                return False
        else:
            return True
    else:
        return False


def is_int(s):
    try: 
        int(s)
        return True
    except ValueError:
        return False
        
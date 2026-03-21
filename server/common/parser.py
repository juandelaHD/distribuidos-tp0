import struct

from common.utils import Bet

"""
Parse a single Bet from its raw wire fields.
bd_bytes must be 4 bytes: year(2, big-endian) | month(1) | day(1).

Raises ValueError if birthdate fields are out of valid calendar range.
"""
def parse_bet(agency, first_name, last_name, document, bd_bytes, number):
    year, month, day = struct.unpack('!HBB', bd_bytes)
    if not (1 <= month <= 12 and 1 <= day <= 31):
        raise ValueError(f"invalid birthdate fields: year={year} month={month} day={day}")
    birthdate = f"{year:04d}-{month:02d}-{day:02d}"
    return Bet(str(agency), first_name, last_name, str(document), birthdate, str(number))

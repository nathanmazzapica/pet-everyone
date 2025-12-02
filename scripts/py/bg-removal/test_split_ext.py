from main import trim_ext
import pytest

@pytest.mark.parametrize(
        "filename, expected",
        [
            ("picture.png", "picture"),
            ("ab3lk23o2392ekadsjdad.jpeg", "ab3lk23o2392ekadsjdad")
        ]
)
def test_trim_ext_valid(filename, expected):
    assert trim_ext(filename) == expected


@pytest.mark.parametrize(
        "filename",
        [
            ("the.picture.png"),
            ("ab3lk23o2392ekadsjdad"),
            (".test"),
            ("."),
            ("")
        ]
)
def test_trim_ext_ValueError(filename):
    with pytest.raises(ValueError):
        trim_ext(filename)



def test_trim_ext_base64():
    import secrets

    for i in range(1000):
        byte_length = 32
        byte_length = 32
        b64 = secrets.token_urlsafe(byte_length)

        encoded_filename = f"{b64}.jpeg"

        assert trim_ext(encoded_filename) == b64




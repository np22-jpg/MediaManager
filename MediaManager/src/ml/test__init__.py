from ml import get_season


def test_get_season():
    for i in range(0, 50):
        assert get_season("The.Big.Bang.Theory.(2007).Season.9.S09.(1080p.BluRay.x265.HEVC.10bit.AAC.5.1.Vyndros)") == 9

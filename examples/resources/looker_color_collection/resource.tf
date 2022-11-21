resource "looker_color_collection" "collection" {
  label      = "My new collection"
  categoricalpalettes {
    label = "cat"
    colors = ["#1A73E8",
        "#12B5CB",
        "#E52592"]
  }
  sequentialpalettes {
    label = "seq"
    stops {
      color = "#FFFFFF"
      offset = "1"
    }
    stops {
      color = "#1A73E8"
      offset = "100"
    }
  }
  divergingpalettes {
    label = "div"
    stops {
      color = "#FFFFFF"
      offset = "1"
    }
    stops {
      color = "#1A73E8"
      offset = "100"
    }
  }
}
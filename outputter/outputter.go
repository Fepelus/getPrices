package outputter

import "github.com/Fepelus/getPrices/entities"

type Outputter interface {
  Append(entities.Price)
  Output()
}


package com.github.kapmahc.axe.nut.repositories;

import com.github.kapmahc.axe.nut.models.Card;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

@Repository("nut.cardRepository")
public interface CardRepository extends CrudRepository<Card, Long> {
}

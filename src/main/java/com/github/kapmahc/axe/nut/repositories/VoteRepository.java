package com.github.kapmahc.axe.nut.repositories;

import com.github.kapmahc.axe.nut.models.Vote;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

@Repository("nut.voteRepository")
public interface VoteRepository  extends CrudRepository<Vote, Long> {
}

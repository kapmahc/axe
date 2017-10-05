package com.github.kapmahc.axe.reading.repositories;

import com.github.kapmahc.axe.reading.models.Note;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

@Repository("ops.reading.NoteRepository")
public interface NoteRepository  extends CrudRepository<Note,Long> {
}

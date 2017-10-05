package com.github.kapmahc.axe.reading.repositories;

import com.github.kapmahc.axe.reading.models.Book;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

@Repository("ops.reading.BookRepository")
public interface BookRepository extends CrudRepository<Book,Long> {
}

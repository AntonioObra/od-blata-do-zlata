//
//  Item.swift
//  od-blata-do-zlata
//
//  Created by Antonio Obradovic on 16.07.2026..
//

import Foundation
import SwiftData

@Model
final class Item {
    var timestamp: Date
    
    init(timestamp: Date) {
        self.timestamp = timestamp
    }
}
